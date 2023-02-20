<?php
/**
 * WorkflowInstanceHistoryStage
 *
 * PHP version 7.4
 *
 * @category Class
 * @package  Formance
 * @author   OpenAPI Generator team
 * @link     https://openapi-generator.tech
 */

/**
 * Formance Stack API
 *
 * Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions />
 *
 * The version of the OpenAPI document: develop
 * Contact: support@formance.com
 * Generated by: https://openapi-generator.tech
 * OpenAPI Generator version: 6.4.0
 */

/**
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

namespace Formance\Model;

use \ArrayAccess;
use \Formance\ObjectSerializer;

/**
 * WorkflowInstanceHistoryStage Class Doc Comment
 *
 * @category Class
 * @package  Formance
 * @author   OpenAPI Generator team
 * @link     https://openapi-generator.tech
 * @implements \ArrayAccess<string, mixed>
 */
class WorkflowInstanceHistoryStage implements ModelInterface, ArrayAccess, \JsonSerializable
{
    public const DISCRIMINATOR = null;

    /**
      * The original name of the model.
      *
      * @var string
      */
    protected static $openAPIModelName = 'WorkflowInstanceHistoryStage';

    /**
      * Array of property to type mappings. Used for (de)serialization
      *
      * @var string[]
      */
    protected static $openAPITypes = [
        'name' => 'string',
        'input' => '\Formance\Model\WorkflowInstanceHistoryStageInput',
        'output' => '\Formance\Model\WorkflowInstanceHistoryStageOutput',
        'error' => 'string',
        'terminated' => 'Bool',
        'started_at' => '\DateTime',
        'terminated_at' => '\DateTime',
        'last_failure' => 'string',
        'attempt' => 'int',
        'next_execution' => '\DateTime'
    ];

    /**
      * Array of property to format mappings. Used for (de)serialization
      *
      * @var string[]
      * @phpstan-var array<string, string|null>
      * @psalm-var array<string, string|null>
      */
    protected static $openAPIFormats = [
        'name' => null,
        'input' => null,
        'output' => null,
        'error' => null,
        'terminated' => null,
        'started_at' => 'date-time',
        'terminated_at' => 'date-time',
        'last_failure' => null,
        'attempt' => null,
        'next_execution' => 'date-time'
    ];

    /**
      * Array of nullable properties. Used for (de)serialization
      *
      * @var boolean[]
      */
    protected static array $openAPINullables = [
        'name' => false,
		'input' => false,
		'output' => false,
		'error' => false,
		'terminated' => false,
		'started_at' => false,
		'terminated_at' => false,
		'last_failure' => false,
		'attempt' => false,
		'next_execution' => false
    ];

    /**
      * If a nullable field gets set to null, insert it here
      *
      * @var boolean[]
      */
    protected array $openAPINullablesSetToNull = [];

    /**
     * Array of property to type mappings. Used for (de)serialization
     *
     * @return array
     */
    public static function openAPITypes()
    {
        return self::$openAPITypes;
    }

    /**
     * Array of property to format mappings. Used for (de)serialization
     *
     * @return array
     */
    public static function openAPIFormats()
    {
        return self::$openAPIFormats;
    }

    /**
     * Array of nullable properties
     *
     * @return array
     */
    protected static function openAPINullables(): array
    {
        return self::$openAPINullables;
    }

    /**
     * Array of nullable field names deliberately set to null
     *
     * @return boolean[]
     */
    private function getOpenAPINullablesSetToNull(): array
    {
        return $this->openAPINullablesSetToNull;
    }

    /**
     * Setter - Array of nullable field names deliberately set to null
     *
     * @param boolean[] $openAPINullablesSetToNull
     */
    private function setOpenAPINullablesSetToNull(array $openAPINullablesSetToNull): void
    {
        $this->openAPINullablesSetToNull = $openAPINullablesSetToNull;
    }

    /**
     * Checks if a property is nullable
     *
     * @param string $property
     * @return bool
     */
    public static function isNullable(string $property): bool
    {
        return self::openAPINullables()[$property] ?? false;
    }

    /**
     * Checks if a nullable property is set to null.
     *
     * @param string $property
     * @return bool
     */
    public function isNullableSetToNull(string $property): bool
    {
        return in_array($property, $this->getOpenAPINullablesSetToNull(), true);
    }

    /**
     * Array of attributes where the key is the local name,
     * and the value is the original name
     *
     * @var string[]
     */
    protected static $attributeMap = [
        'name' => 'name',
        'input' => 'input',
        'output' => 'output',
        'error' => 'error',
        'terminated' => 'terminated',
        'started_at' => 'startedAt',
        'terminated_at' => 'terminatedAt',
        'last_failure' => 'lastFailure',
        'attempt' => 'attempt',
        'next_execution' => 'nextExecution'
    ];

    /**
     * Array of attributes to setter functions (for deserialization of responses)
     *
     * @var string[]
     */
    protected static $setters = [
        'name' => 'setName',
        'input' => 'setInput',
        'output' => 'setOutput',
        'error' => 'setError',
        'terminated' => 'setTerminated',
        'started_at' => 'setStartedAt',
        'terminated_at' => 'setTerminatedAt',
        'last_failure' => 'setLastFailure',
        'attempt' => 'setAttempt',
        'next_execution' => 'setNextExecution'
    ];

    /**
     * Array of attributes to getter functions (for serialization of requests)
     *
     * @var string[]
     */
    protected static $getters = [
        'name' => 'getName',
        'input' => 'getInput',
        'output' => 'getOutput',
        'error' => 'getError',
        'terminated' => 'getTerminated',
        'started_at' => 'getStartedAt',
        'terminated_at' => 'getTerminatedAt',
        'last_failure' => 'getLastFailure',
        'attempt' => 'getAttempt',
        'next_execution' => 'getNextExecution'
    ];

    /**
     * Array of attributes where the key is the local name,
     * and the value is the original name
     *
     * @return array
     */
    public static function attributeMap()
    {
        return self::$attributeMap;
    }

    /**
     * Array of attributes to setter functions (for deserialization of responses)
     *
     * @return array
     */
    public static function setters()
    {
        return self::$setters;
    }

    /**
     * Array of attributes to getter functions (for serialization of requests)
     *
     * @return array
     */
    public static function getters()
    {
        return self::$getters;
    }

    /**
     * The original name of the model.
     *
     * @return string
     */
    public function getModelName()
    {
        return self::$openAPIModelName;
    }


    /**
     * Associative array for storing property values
     *
     * @var mixed[]
     */
    protected $container = [];

    /**
     * Constructor
     *
     * @param mixed[] $data Associated array of property values
     *                      initializing the model
     */
    public function __construct(array $data = null)
    {
        $this->setIfExists('name', $data ?? [], null);
        $this->setIfExists('input', $data ?? [], null);
        $this->setIfExists('output', $data ?? [], null);
        $this->setIfExists('error', $data ?? [], null);
        $this->setIfExists('terminated', $data ?? [], null);
        $this->setIfExists('started_at', $data ?? [], null);
        $this->setIfExists('terminated_at', $data ?? [], null);
        $this->setIfExists('last_failure', $data ?? [], null);
        $this->setIfExists('attempt', $data ?? [], null);
        $this->setIfExists('next_execution', $data ?? [], null);
    }

    /**
    * Sets $this->container[$variableName] to the given data or to the given default Value; if $variableName
    * is nullable and its value is set to null in the $fields array, then mark it as "set to null" in the
    * $this->openAPINullablesSetToNull array
    *
    * @param string $variableName
    * @param array  $fields
    * @param mixed  $defaultValue
    */
    private function setIfExists(string $variableName, array $fields, $defaultValue): void
    {
        if (self::isNullable($variableName) && array_key_exists($variableName, $fields) && is_null($fields[$variableName])) {
            $this->openAPINullablesSetToNull[] = $variableName;
        }

        $this->container[$variableName] = $fields[$variableName] ?? $defaultValue;
    }

    /**
     * Show all the invalid properties with reasons.
     *
     * @return array invalid properties with reasons
     */
    public function listInvalidProperties()
    {
        $invalidProperties = [];

        if ($this->container['name'] === null) {
            $invalidProperties[] = "'name' can't be null";
        }
        if ($this->container['input'] === null) {
            $invalidProperties[] = "'input' can't be null";
        }
        if ($this->container['terminated'] === null) {
            $invalidProperties[] = "'terminated' can't be null";
        }
        if ($this->container['started_at'] === null) {
            $invalidProperties[] = "'started_at' can't be null";
        }
        if ($this->container['attempt'] === null) {
            $invalidProperties[] = "'attempt' can't be null";
        }
        return $invalidProperties;
    }

    /**
     * Validate all the properties in the model
     * return true if all passed
     *
     * @return bool True if all properties are valid
     */
    public function valid()
    {
        return count($this->listInvalidProperties()) === 0;
    }


    /**
     * Gets name
     *
     * @return string
     */
    public function getName()
    {
        return $this->container['name'];
    }

    /**
     * Sets name
     *
     * @param string $name name
     *
     * @return self
     */
    public function setName($name)
    {
        if (is_null($name)) {
            throw new \InvalidArgumentException('non-nullable name cannot be null');
        }
        $this->container['name'] = $name;

        return $this;
    }

    /**
     * Gets input
     *
     * @return \Formance\Model\WorkflowInstanceHistoryStageInput
     */
    public function getInput()
    {
        return $this->container['input'];
    }

    /**
     * Sets input
     *
     * @param \Formance\Model\WorkflowInstanceHistoryStageInput $input input
     *
     * @return self
     */
    public function setInput($input)
    {
        if (is_null($input)) {
            throw new \InvalidArgumentException('non-nullable input cannot be null');
        }
        $this->container['input'] = $input;

        return $this;
    }

    /**
     * Gets output
     *
     * @return \Formance\Model\WorkflowInstanceHistoryStageOutput|null
     */
    public function getOutput()
    {
        return $this->container['output'];
    }

    /**
     * Sets output
     *
     * @param \Formance\Model\WorkflowInstanceHistoryStageOutput|null $output output
     *
     * @return self
     */
    public function setOutput($output)
    {
        if (is_null($output)) {
            throw new \InvalidArgumentException('non-nullable output cannot be null');
        }
        $this->container['output'] = $output;

        return $this;
    }

    /**
     * Gets error
     *
     * @return string|null
     */
    public function getError()
    {
        return $this->container['error'];
    }

    /**
     * Sets error
     *
     * @param string|null $error error
     *
     * @return self
     */
    public function setError($error)
    {
        if (is_null($error)) {
            throw new \InvalidArgumentException('non-nullable error cannot be null');
        }
        $this->container['error'] = $error;

        return $this;
    }

    /**
     * Gets terminated
     *
     * @return Bool
     */
    public function getTerminated()
    {
        return $this->container['terminated'];
    }

    /**
     * Sets terminated
     *
     * @param Bool $terminated terminated
     *
     * @return self
     */
    public function setTerminated($terminated)
    {
        if (is_null($terminated)) {
            throw new \InvalidArgumentException('non-nullable terminated cannot be null');
        }
        $this->container['terminated'] = $terminated;

        return $this;
    }

    /**
     * Gets started_at
     *
     * @return \DateTime
     */
    public function getStartedAt()
    {
        return $this->container['started_at'];
    }

    /**
     * Sets started_at
     *
     * @param \DateTime $started_at started_at
     *
     * @return self
     */
    public function setStartedAt($started_at)
    {
        if (is_null($started_at)) {
            throw new \InvalidArgumentException('non-nullable started_at cannot be null');
        }
        $this->container['started_at'] = $started_at;

        return $this;
    }

    /**
     * Gets terminated_at
     *
     * @return \DateTime|null
     */
    public function getTerminatedAt()
    {
        return $this->container['terminated_at'];
    }

    /**
     * Sets terminated_at
     *
     * @param \DateTime|null $terminated_at terminated_at
     *
     * @return self
     */
    public function setTerminatedAt($terminated_at)
    {
        if (is_null($terminated_at)) {
            throw new \InvalidArgumentException('non-nullable terminated_at cannot be null');
        }
        $this->container['terminated_at'] = $terminated_at;

        return $this;
    }

    /**
     * Gets last_failure
     *
     * @return string|null
     */
    public function getLastFailure()
    {
        return $this->container['last_failure'];
    }

    /**
     * Sets last_failure
     *
     * @param string|null $last_failure last_failure
     *
     * @return self
     */
    public function setLastFailure($last_failure)
    {
        if (is_null($last_failure)) {
            throw new \InvalidArgumentException('non-nullable last_failure cannot be null');
        }
        $this->container['last_failure'] = $last_failure;

        return $this;
    }

    /**
     * Gets attempt
     *
     * @return int
     */
    public function getAttempt()
    {
        return $this->container['attempt'];
    }

    /**
     * Sets attempt
     *
     * @param int $attempt attempt
     *
     * @return self
     */
    public function setAttempt($attempt)
    {
        if (is_null($attempt)) {
            throw new \InvalidArgumentException('non-nullable attempt cannot be null');
        }
        $this->container['attempt'] = $attempt;

        return $this;
    }

    /**
     * Gets next_execution
     *
     * @return \DateTime|null
     */
    public function getNextExecution()
    {
        return $this->container['next_execution'];
    }

    /**
     * Sets next_execution
     *
     * @param \DateTime|null $next_execution next_execution
     *
     * @return self
     */
    public function setNextExecution($next_execution)
    {
        if (is_null($next_execution)) {
            throw new \InvalidArgumentException('non-nullable next_execution cannot be null');
        }
        $this->container['next_execution'] = $next_execution;

        return $this;
    }
    /**
     * Returns true if offset exists. False otherwise.
     *
     * @param integer $offset Offset
     *
     * @return boolean
     */
    public function offsetExists($offset): bool
    {
        return isset($this->container[$offset]);
    }

    /**
     * Gets offset.
     *
     * @param integer $offset Offset
     *
     * @return mixed|null
     */
    #[\ReturnTypeWillChange]
    public function offsetGet($offset)
    {
        return $this->container[$offset] ?? null;
    }

    /**
     * Sets value based on offset.
     *
     * @param int|null $offset Offset
     * @param mixed    $value  Value to be set
     *
     * @return void
     */
    public function offsetSet($offset, $value): void
    {
        if (is_null($offset)) {
            $this->container[] = $value;
        } else {
            $this->container[$offset] = $value;
        }
    }

    /**
     * Unsets offset.
     *
     * @param integer $offset Offset
     *
     * @return void
     */
    public function offsetUnset($offset): void
    {
        unset($this->container[$offset]);
    }

    /**
     * Serializes the object to a value that can be serialized natively by json_encode().
     * @link https://www.php.net/manual/en/jsonserializable.jsonserialize.php
     *
     * @return mixed Returns data which can be serialized by json_encode(), which is a value
     * of any type other than a resource.
     */
    #[\ReturnTypeWillChange]
    public function jsonSerialize()
    {
       return ObjectSerializer::sanitizeForSerialization($this);
    }

    /**
     * Gets the string presentation of the object
     *
     * @return string
     */
    public function __toString()
    {
        return json_encode(
            ObjectSerializer::sanitizeForSerialization($this),
            JSON_PRETTY_PRINT
        );
    }

    /**
     * Gets a header-safe presentation of the object
     *
     * @return string
     */
    public function toHeaderValue()
    {
        return json_encode(ObjectSerializer::sanitizeForSerialization($this));
    }
}


