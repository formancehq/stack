<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


/**
 * CreateBalanceResponse - Created balance
 * 
 * @package formance\stack\Models\Shared
 * @access public
 */
class CreateBalanceResponse
{
	#[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\WalletsBalance')]
    public WalletsBalance $data;
    
	public function __construct()
	{
		$this->data = new \formance\stack\Models\Shared\WalletsBalance();
	}
}
